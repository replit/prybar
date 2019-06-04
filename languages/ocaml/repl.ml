#use "topfind";;
#require "compiler-libs";;
#require "unix";;
(*#require "reason";;*)
#load "ocamlcommon.cma";;
#load "unix.cma";;
(*#load "reason.cma";;*)

type syntax = Reason | OCaml

type eval_result = Success | Error of string

type mode =
  | RunCode of string * syntax
  | StartRepl of syntax
  | PrintExpression of string * syntax
  | Stdin of syntax
  | Invalid of string

type args = {mode: mode}

exception LangException of string

let _ =
  let interactive_mode = ref false in
  let eval_filepath = ref None in
  let quiet = ref false in
  let syntax = ref OCaml in
  let print_mode = ref false in
  let run_mode = ref false in
  let ps1 = Sys.getenv_opt "PRYBAR_PS1" |> function | Some str -> str | None -> "#" in
  let ps2 = Sys.getenv_opt "PRYBAR_PS2" |> function | Some str -> str | None -> " " in
  (* If there is code provided to interpret *)
  let code = ref "" in
  let parse_syntax = function
    | "re" -> syntax := Reason
    | "ml" -> syntax := OCaml
    | arg -> raise (Arg.Bad ("Unknown syntax: " ^ arg))
  in
  let print_mode_arg str =
    print_mode := true ;
    code := str
  in
  let run_mode_arg str =
    run_mode := true ;
    code := str
  in
  let speclist =
    [ ("-q", Arg.Set quiet, "Don't print OCaml version on startup")
    ; ( "-e" , Arg.String print_mode_arg , "Eval and output results of interpreted code" )
    ; ("-c", Arg.String run_mode_arg, "Run code without printing the output")
    ; ("-i", Arg.Set interactive_mode, "Run as interactive repl")
    ; ( "-s" , Arg.Symbol (["ml"; "re"], parse_syntax), " Sets the syntax for the repl (default: ml)" ) ]
  in
  let usage_msg =
    "OCaml / Reason repl script for prybar. Options available:"
  in
  Arg.parse speclist
    (fun str ->
      match (str, !interactive_mode) with
      | (path, true) -> eval_filepath := Some(path)
      | _ -> print_endline ("Anonymous arg: " ^ str)
      )
    usage_msg ;
  let mode =
    match (!print_mode, !run_mode, !interactive_mode) with
    | true, false, false -> PrintExpression (!code, !syntax)
    | false, true, false -> RunCode (!code, !syntax)
    | false, false, true -> StartRepl !syntax
    | _ -> Invalid "You can only use one mode: -e | -i | -c"
  in
  let module Custom_Toploop = struct
    include Toploop
    open Format

    (* Most stuff is straightly copied and modified
     * from github.com/ocaml/ocaml/blob/trunk/toplevel/toploop.ml *)
    exception PPerror

    let first_line = ref true

    let got_eof = ref false

    let refill_lexbuf buffer len =
      if !got_eof then (
        got_eof := false ;
        0 )
      else
        let prompt =
          if !first_line then ps1 ^ " "
          else if Lexer.in_comment () then "* "
          else ps2 (* continuation prompt *)
        in
        first_line := false ;
        let len, eof = !read_interactive_input prompt buffer len in
        if eof then (
          Location.echo_eof () ;
          if len > 0 then got_eof := true ;
          len )
        else len

    (* Minimal version of just running any input file, we stripped a lot of original logic
     * because we don't want to do any side effects on the compiler environment *)
    let run_script ppf name =
      let explicit_name =
        (* Prevent use_silently from searching in the path. *)
        if name <> "" && Filename.is_implicit name
        then Filename.concat Filename.current_dir_name name
        else name
      in
      use_silently ppf explicit_name

    let loop ppf =
      Clflags.debug := true ;
      Location.formatter_for_warnings := ppf ;
      ( try initialize_toplevel_env () with
      | (Env.Error _ | Typetexp.Error _) as exn ->
          Location.report_exception ppf exn ;
          exit 2 ) ;
      let lb = Lexing.from_function refill_lexbuf in
      Location.init lb "//toplevel//" ;
      Location.input_name := "//toplevel//" ;
      Location.input_lexbuf := Some lb ;
      Sys.catch_break true ;
      (*load_ocamlinit ppf;*)

      (* If there's an entry file provided, run it before dropping into interactive mode *)
      (match !eval_filepath with
      | Some name -> run_script ppf name
      | _ -> false) |> ignore ;

      while true do
        let snap = Btype.snapshot () in
        try
          Lexing.flush_input lb ;
          Location.reset () ;
          Warnings.reset_fatal () ;
          first_line := true ;
          let phr =
            try !parse_toplevel_phrase lb with Exit -> raise PPerror
          in
          let phr = preprocess_phrase ppf phr in
          Env.reset_cache_toplevel () ;
          ignore (execute_phrase true ppf phr)
        with
        | End_of_file -> exit 0
        | Sys.Break ->
            fprintf ppf "Interrupted.@." ;
            Btype.backtrack snap
        | PPerror -> ()
        | x ->
            Location.report_exception ppf x ;
            Btype.backtrack snap
      done
  end in
  let module Repl = struct
    let std_fmt = Format.std_formatter

    let noop_fmt = Format.make_formatter (fun _ _ _ -> ()) ignore

    let init_toploop () =
      Topfind.add_predicates ["byte"; "toploop"] ;
      (* Add findlib path so Topfind is available and it won't be
           initialized twice if the user does [#use "topfind"]. *)
      Topdirs.dir_directory (Findlib.package_directory "findlib") ;
      Toploop.initialize_toplevel_env ()

    let start_loop ?(fmt = std_fmt) () = Custom_Toploop.loop fmt

    let eval ?(fmt = std_fmt) ~syntax str =
      try
        let lex = Lexing.from_string str in
        let tpl_phrases =
          match syntax with
          | OCaml | Reason -> Parse.use_file lex
          (*| Reason ->*)
              (*List.map Reason_toolchain.To_current.copy_toplevel_phrase*)
                (*(Reason_toolchain.RE.use_file lex)*)
        in
        let exec phr =
          if Toploop.execute_phrase true fmt phr then Success
          else Error "No result"
        in
        let rec execAll phrases =
          match phrases with
          | [] -> Error "No result"
          | [phr] -> exec phr
          | phr :: next -> (
              let ret = exec phr in
              match ret with Error _ -> ret | _ -> execAll next )
        in
        execAll tpl_phrases
      with
      | Syntaxerr.Error _ -> Error "Syntax Error occurred"
      (*| Reason_syntax_util.Error _ -> Error "Reason Parsing Error"*)
      | _ -> Error ("Error while exec: " ^ str)
  end in
  if not !quiet then print_endline ("OCaml " ^ Sys.ocaml_version ^ " on " ^ Sys.os_type);
  match mode with
  | PrintExpression (code, syntax) ->
      Repl.init_toploop () ;
      Repl.eval ~syntax code |> ignore
  | RunCode (code, syntax) ->
      Repl.init_toploop () ;
      Repl.eval ~syntax ~fmt:Repl.noop_fmt code |> ignore
  | StartRepl syntax ->
    Repl.init_toploop () ;
    Repl.start_loop ()
  | Stdin syntax -> print_endline "Reading from stdin"
  | Invalid str -> print_endline str
