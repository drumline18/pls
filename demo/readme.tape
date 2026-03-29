Output assets/pls-demo.webm

Set Shell "bash"
Set Width 1200
Set Height 720
Set FontSize 28
Set TypingSpeed 35ms
Set WaitTimeout 20s

# Show a low-risk listing suggestion.
Type "./bin/pls --no-exec list all dotfiles in the current directory"
Enter
Wait+Screen /Risk:/
Sleep 1500ms

Type "clear"
Enter
Sleep 300ms

# Show a concise service-check suggestion.
Type "./bin/pls --no-exec check if jellyfin is running"
Enter
Wait+Screen /Risk:/
Sleep 1700ms
