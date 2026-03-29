Output assets/pls-demo.webm

Set Shell "bash"
Set Width 1200
Set Height 720
Set FontSize 28
Set TypingSpeed 35ms
Set WaitTimeout 20s

# Put the freshly built binary on PATH so the demo can use plain `pls`.
Type "export PATH=$PWD/bin:$PATH"
Enter
Sleep 400ms
Type "clear"
Enter
Sleep 300ms

# Show natural language -> command -> confirmation.
Type "pls list all dotfiles in the current directory"
Enter
Wait+Screen /Run it\?/
Sleep 1500ms
Type "n"
Enter
Sleep 700ms

Type "clear"
Enter
Sleep 300ms

# Show a second plain-language example with the same confirmation flow.
Type "pls check if jellyfin is running"
Enter
Wait+Screen /Run it\?/
Sleep 1700ms
Type "n"
Enter
Sleep 900ms
