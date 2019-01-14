type mode = Execute of string | Repl | Invalid

let _ =
  Sys.argv |> Array.iteri (fun i x -> print_endline ((string_of_int i) ^ ":" ^ x)) |> ignore;
  let mode = match Sys.argv with
    | [|_; "-e"; code|] -> Execute(code)
    | [|_; "-i"|] -> Repl 
    | _ -> Invalid

  in
  match mode with
    | Execute code -> print_endline ("Would execute: " ^ code)
    | Repl -> print_endline "Would enter Repl mode"
    | Invalid -> print_endline "Invalid mode";;

