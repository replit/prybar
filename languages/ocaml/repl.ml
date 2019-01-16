#use "topfind";;
#require "compiler-libs";;
#require "unix";;
#require "reason";;
#load "ocamlcommon.cma";;
#load "unix.cma";;
#load "reason.cma";;

type syntax = Reason | OCaml

type eval_result = Success | Error of string

type mode =
  | RunCode of string * syntax
  | StartRepl of syntax
  | PrintExpression of string * syntax
  | Stdin of syntax
  | Invalid

exception LangException of string

let _ =
  let parseSyntax = function
    | "reason" -> Reason
    | "ocaml" -> OCaml
    | arg -> raise (LangException ("Unknown syntax: " ^ arg))
  in
  let mode =
    match Sys.argv with
    | [|_; syntax; "-e"; code|] -> PrintExpression (code, parseSyntax syntax)
    | [|_; syntax; "-c"; code|] -> RunCode (code, parseSyntax syntax)
    | [|_; syntax; "-stdin"|] -> Stdin (parseSyntax syntax)
    | [|_; syntax; "-i"|] -> StartRepl (parseSyntax syntax)
    | _ -> Invalid
  in
  let module Repl = struct
    let std_fmt = Format.std_formatter
    let noop_fmt = Format.make_formatter (fun _ _ _ -> ()) ignore

    let init_toploop () =
      Topfind.add_predicates ["byte"; "toploop"] ;
      (* Add findlib path so Topfind is available and it won't be
           initialized twice if the user does [#use "topfind"]. *)
      Topdirs.dir_directory (Findlib.package_directory "findlib") ;
      Toploop.initialize_toplevel_env ()

    let eval ?(fmt = std_fmt) ~syntax str =
      try
        let lex = Lexing.from_string str in
        let tpl_phrases =
          match syntax with
          | OCaml -> Parse.use_file lex
          | Reason ->
              List.map Reason_toolchain.To_current.copy_toplevel_phrase
                (Reason_toolchain.RE.use_file lex)
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
      | Reason_syntax_util.Error _ -> Error "Reason Parsing Error"
      | _ -> Error ("Error while exec: " ^ str)
  end in
  match mode with
  | PrintExpression (code, syntax) -> Repl.init_toploop () ; Repl.eval ~syntax code |> ignore
  | RunCode (code, syntax) -> Repl.init_toploop () ; Repl.eval ~syntax ~fmt:Repl.noop_fmt code |> ignore
  | StartRepl syntax -> print_endline "Would enter Repl mode"
  | Stdin syntax -> print_endline "Reading from stdin"
  | Invalid -> print_endline "Invalid mode"
