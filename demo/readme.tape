Output assets/pls-demo.webm

Set Shell "bash"
Set Width 1200
Set Height 720
Set FontSize 28
Set TypingSpeed 35ms
Set WaitTimeout 20s

# Hidden bootstrap so the visible demo starts clean.
Hide
Type "mkdir -p $HOME/.cache/pls-vhs-demo"
Enter
Type "touch $HOME/.cache/pls-vhs-demo/.alpha $HOME/.cache/pls-vhs-demo/.beta"
Enter
Type "cd $HOME/.cache/pls-vhs-demo"
Enter
Type "clear"
Enter
Show
Sleep 300ms

# Show natural language -> command -> confirmation -> real execution.
Type "pls list all dotfiles in the current directory"
Enter
Wait+Screen /Run it\?/
Sleep 1200ms
Type "y"
Enter
Sleep 2200ms

Type "clear"
Enter
Sleep 300ms

# Show a second plain-language example that executes and returns a real status.
Type "pls check if jellyfin is running"
Enter
Wait+Screen /Run it\?/
Sleep 1400ms
Type "y"
Enter
Sleep 2400ms
