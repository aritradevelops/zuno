#!/bin/bash

# Replace <your-username> with your GitHub username or organization name
USERNAME="aritradevelops"

# Add --limit <number> to change more than the default 30 repositories
gh repo list $USERNAME --limit 100 --json name --jq '.[].name' | while IFS= read -r repo
do
    gh repo edit "$USERNAME/$repo" --visibility private --accept-visibility-change-consequences
done

