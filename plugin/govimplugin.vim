let s:channel = ""
let s:timer = ""
let s:plugindir = "/home/tom/src/govimplugin"

let opts = {"in_mode": "json", "out_mode": "json", "err_mode": "json", "callback": function("s:define"), "timeout": 30000}
if $GOVIMTEST_SOCKET != ""
  let s:channel = ch_open($GOVIMTEST_SOCKET, opts)
else
  let targetbin = s:install()
  let opts.exit_cb = function("s:govimExit")
  let job = job_start(targetbin, opts)
  let s:channel = job_getchannel(job)
	echomsg "govimplugin started"
endif

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
	echomsg "receiving ".l:msg
  try
    let l:id = a:msg[0]
    let l:resp = ["callback", l:id, [""]]
    if a:msg[1] == "loaded"
      let s:plugin_status = "loaded"
      "for F in s:loadStatusCallbacks
      "  call call(F, [s:govim_status])
      "endfor
    elseif a:msg[1] == "???" " extend msg handling?
		endif
	catch
		let l:resp[2][0] = 'Caught ' . string(v:exception) . ' in ' . v:throwpoint
  endtry
  call ch_sendexpr(a:channel, l:resp)
endfunc
