param($musicFile, $timbre)

if (-Not($musicFile)) {
    Write-Host "Usage :"$MyInvocation.MyCommand.Name"musicDataFile <timbre>"
    exit
}

if (-Not($timbre)) {
    $timbre = 1     # piano
}

go run Go_PlayBox.go $musicFile $timbre
