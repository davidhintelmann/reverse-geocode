# Configuration variables for pathing to system variable,
# module directory, and go file. 
$WorkDir = "C:\Users\david\go\src\github.com\davidhintelmann\reverse-geocode\cmd"
$EmbedDir = "C:\Users\david\go\src\github.com\davidhintelmann\reverse-geocode\node"
$SystemVariable = "go"
$EmbedFile = "embed_csv.go"
$File = "main.go"
$Binary = ".\geo.exe"
$ArgumentBuild1 = ("run " + $EmbedFile)
$ArgumentBuild2 = ("build -o " + $Binary + " " + $File)

Write-Host "embedding go binary..."
# Effectively running 'go build -o main.exe main.go' in terminal/PowerShell
# but instead using PowerShell functions to measure execution duration.
$BuildTime = Measure-Command {
    Start-Process -FilePath $SystemVariable -WorkingDirectory $EmbedDir -ArgumentList $ArgumentBuild1 -NoNewWindow -wait
}
Write-Host ("Time taken to embed go binary -> " + $BuildTime)


Write-Host "building go binary..."
$BuildTime = Measure-Command {
    Start-Process -FilePath $SystemVariable -WorkingDirectory $WorkDir -ArgumentList $ArgumentBuild2 -NoNewWindow -wait
}
Write-Host ("Time taken to build go binary -> " + $BuildTime)

# Now run the go binary
$RunBinaryTime = Measure-Command {
    Start-Process -FilePath $Binary -WorkingDirectory $WorkDir -NoNewWindow -wait
}
Write-Host ("Time taken to run go binary -> " + $RunBinaryTime)