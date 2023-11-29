# vim: ft=kak

define-command luapls-log %{
  edit /tmp/luapls.log -readonly
}

hook global WinCreate /tmp/luapls\.log %{
  add-highlighter window/luapls-datetime regex '\d*/\d*/\d* \d*:\d*:\d*\.\d*' 0:comment
  add-highlighter window/luapls-server regex ' (\[[\w\.]*?\]) ' 1:value
  add-highlighter window/luapls-rpc regex ' (\[luapls\.rpc\]) ' 1:comment
  add-highlighter window/luapls-error regex '[^\[]*?(ERROR)' 1:Error
  add-highlighter window/luapls-info regex '[^\[]*?(INFO)' 1:function
  add-highlighter window/luapls-notice regex '[^\[]*?(NOTICE)' 1:string
  add-highlighter window/luapls-warning regex '[^\[]*?(WARN)' 1:type
  add-highlighter window/luapls-critical regex '[^\[]*?(CRITICAL)' 1:Error
  add-highlighter window/luapls-server-send regex ' (<--) ' 1:function
  add-highlighter window/luapls-server-recv regex ' (-->) ' 1:string

  evaluate-commands -draft %{
    try %{
      execute-keys '%sreading from stdin, writing to stdout\n<ret>x'
      evaluate-commands -itersel %{
        add-highlighter window/ line %val{cursor_line} +r@keyword
      }
    }
  }
}

hook global BufReload .*luapls\.log %{
  execute-keys 'gjgh'
}
