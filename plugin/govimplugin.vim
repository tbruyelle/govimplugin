let g:channel = "" " TODO remove global scope
let s:timer = ""
let s:plugindir = "/home/tom/src/govimplugin"


function s:callbackCommand(name, flags, ...)
  let l:args = ["function", "command:".a:name, a:flags]
  call extend(l:args, a:000)
	let l:resp = ch_evalexpr(g:channel, l:args)
	if l:resp[0] != "" 
		throw l:resp[0]
	endif
	return l:resp[1]
endfunction

func s:defineCommand(name, attrs)
  let l:def = "command! "
  let l:args = ""
  let l:flags = ['"mods": expand("<mods>")']
  " let l:flags = []
  if has_key(a:attrs, "nargs")
    let l:def .= " ". a:attrs["nargs"]
    if a:attrs["nargs"] != "-nargs=0"
      let l:args = ", <f-args>"
    endif
  endif
  if has_key(a:attrs, "range")
    let l:def .= " ".a:attrs["range"]
    call add(l:flags, '"line1": <line1>')
    call add(l:flags, '"line2": <line2>')
    call add(l:flags, '"range": <range>')
	endif
	if has_key(a:attrs, "count")
    let l:def .= " ". a:attrs["count"]
    call add(l:flags, '"count": <count>')
  endif
  if has_key(a:attrs, "complete")
    let l:def .= " ". a:attrs["complete"]
  endif
  if has_key(a:attrs, "general")
    for l:a in a:attrs["general"]
      let l:def .= " ". l:a
      if l:a == "-bang"
        call add(l:flags, '"bang": "<bang>"')
      endif
      if l:a == "-register"
        call add(l:flags, '"register": "<reg>"')
      endif
    endfor
  endif
  let l:flagsStr = "{" . join(l:flags, ", ") . "}"
  let l:def .= " " . a:name . " call s:callbackCommand(\"". a:name . "\", " . l:flagsStr . l:args . ")"
  execute l:def
endfunction

func s:install()
	let oldpath = getcwd()
  execute "cd ".s:plugindir


	let install = system("go install .")
	if v:shell_error
		throw install
	endif
	execute "cd ".oldpath
	return "/home/tom/go/bin/govimplugin"
endfunc

func s:define(channel, msg)
	" format is [id, type, ...]
  " type is function, command or autocmd
  try
    let l:id = a:msg[0]
    let l:resp = ["callback", l:id, [""]]
    if a:msg[1] == "loaded"
      let s:plugin_status = "loaded"
      "for F in s:loadStatusCallbacks
      "  call call(F, [s:govim_status])
      "endfor
    elseif a:msg[1] == "command"
      call s:defineCommand(a:msg[2], a:msg[3])
    elseif a:msg[1] == "???" " extend msg handling?
		endif
	catch
		let l:resp[2][0] = 'Caught ' . string(v:exception) . ' in ' . v:throwpoint
  endtry
  call ch_sendexpr(a:channel, l:resp)
endfunc

function s:pluginExit(job, exitstatus)
  if a:exitstatus != 0
    let s:govim_status = "failed"
  else
    let s:govim_status = "exited"
  endif
  "for i in s:loadStatusCallbacks
  "  call call(i, [s:govim_status])
  "endfor
  if a:exitstatus != 0
    throw "govim plugin died :("
  endif
endfunction

let opts = {"in_mode": "json", "out_mode": "json", "err_mode": "json", "callback": function("s:define"), "timeout": 30000}
if $GOVIMTEST_SOCKET != ""
  let g:channel = ch_open($GOVIMTEST_SOCKET, opts)
else
  let targetbin = s:install()
  let opts.exit_cb = function("s:pluginExit")
  let job = job_start(targetbin, opts)
  let g:channel = job_getchannel(job)
endif
