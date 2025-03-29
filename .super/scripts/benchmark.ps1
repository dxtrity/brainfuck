Measure-Command { 
    build/brainfuck.exe  bf/bad.bf | Out-Default
} | Format-List Milliseconds