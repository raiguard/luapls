# vim: ft=kak

hook global WinCreate \*debug\* %{
  add-highlighter window/luapls-newserver regex 'reading from stdin, writing to stdout' 0:+r@keyword
  add-highlighter window/luapls-datetime regex '\d*/\d*/\d* \d*:\d*:\d*\.\d*' 0:comment
  add-highlighter window/luapls-server regex ' (\[[\w\.]*?\]) ' 1:value
  add-highlighter window/luapls-rpc regex ' (\[luapls\.rpc\]) ' 1:comment
  add-highlighter window/luapls-error regex '[^\[]*?(ERROR)' 1:Error
  add-highlighter window/luapls-info regex '[^\[]*?(INFO)' 1:function
  add-highlighter window/luapls-notice regex '[^\[]*?(NOTE)' 1:string
  add-highlighter window/luapls-warning regex '[^\[]*?(WARN)' 1:type
  add-highlighter window/luapls-critical regex '[^\[]*?(CRIT)' 1:+r@Error
  add-highlighter window/luapls-server-send regex ' (<--) ' 1:function
  add-highlighter window/luapls-server-recv regex ' (-->) ' 1:string
}
