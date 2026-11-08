Set-Location "c:\Users\11\Desktop\github proj\sharedrop"

$env:GIT_AUTHOR_NAME    = "nw4m"
$env:GIT_AUTHOR_EMAIL   = "239559894+nw4m@users.noreply.github.com"
$env:GIT_COMMITTER_NAME  = "nw4m"
$env:GIT_COMMITTER_EMAIL = "239559894+nw4m@users.noreply.github.com"
$env:FILTER_BRANCH_SQUELCH_WARNING = "1"

# new messages - natural, lazy human style
$messages = @{
    "384ed62dcc8eed04bd7f5bbd340927db8a58f175" = "initial commit"
    "48d2537be05efbe1f1818051d7ecaf83af993114" = "add core server and cli"
    "12fab19692c30f605f6b3a4dffeb4a9418597b74" = "add install command"
    "c8f51e2a8f55c2a0654eb9d501b6028967e6cca5" = "cleanup file icons"
    "ab6beb6f3513071db5645685d4cd9f3404afc559" = "new ui + root flag"
    "c1832fa5d08e870436a85a71c8b0275caa1b22fd" = "fix upload"
    "f9423c057a190330c12013a0a4983fc89aff241e" = "dark mode, qr code, progress bar"
}

# use git filter-branch with msg filter
$scriptBlock = @"
case `$GIT_COMMIT in
    384ed62dcc8eed04bd7f5bbd340927db8a58f175) echo "initial commit";;
    48d2537be05efbe1f1818051d7ecaf83af993114) echo "add core server and cli";;
    12fab19692c30f605f6b3a4dffeb4a9418597b74) echo "add install command";;
    c8f51e2a8f55c2a0654eb9d501b6028967e6cca5) echo "cleanup file icons";;
    ab6beb6f3513071db5645685d4cd9f3404afc559) echo "new ui + root flag";;
    c1832fa5d08e870436a85a71c8b0275caa1b22fd) echo "fix upload";;
    f9423c057a190330c12013a0a4983fc89aff241e) echo "dark mode, qr code, progress bar";;
    *) cat;;
esac
"@

# write the script to a temp file
$tmpScript = "$env:TEMP\git_msg_filter.sh"
$scriptBlock | Out-File -FilePath $tmpScript -Encoding UTF8 -NoNewline

Write-Host "Script written to $tmpScript"
Write-Host "Running filter-branch..."

git filter-branch -f `
    --msg-filter "sh '$($tmpScript -replace '\\','/')'" `
    --env-filter '
        GIT_AUTHOR_NAME="nw4m"
        GIT_AUTHOR_EMAIL="239559894+nw4m@users.noreply.github.com"
        GIT_COMMITTER_NAME="nw4m"
        GIT_COMMITTER_EMAIL="239559894+nw4m@users.noreply.github.com"
    ' `
    -- --all

Write-Host "REWRITE_DONE:$LASTEXITCODE"
