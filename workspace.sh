#! /bin/bash

session='ai-solver'
session_exists=$(tmux list-sessions | grep $session)

if [ "$session_exists" = "" ]
then
    tmux new-session -d -s $session

    # Main window for GIT
    tmux rename-window -t 0 'main'

    # API window
    tmux new-window -t $session:1 -n 'api'
    tmux send-keys -t 'api' 'cd api' C-m 'go run .' C-m

    # Frontend window
    tmux new-window -t $session:2 -n 'frontend'
    tmux send-keys -t 'frontend' 'cd frontend' C-m 'npm run dev' C-m

    # Frontend window
    tmux new-window -t $session:3 -n 'misc'
fi

tmux attach-session -t $session:0