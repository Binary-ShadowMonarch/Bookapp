#!/bin/bash

SESSION_NAME="bookapp-logs"

# Kill the session if it already exists to ensure a clean start
tmux has-session -t $SESSION_NAME 2>/dev/null
if [ $? = 0 ]; then
    tmux kill-session -t $SESSION_NAME
fi

# Create a new detached session with a window named "Docker"
# The initial pane is 0
tmux new-session -d -s $SESSION_NAME -n "Docker"

# Create the 2x2 layout
# Split pane 0 vertically, creating pane 1 below it
tmux split-window -v -t $SESSION_NAME:Docker.0
# Split the top pane (0) horizontally, creating pane 2 to its right
tmux split-window -h -t $SESSION_NAME:Docker.0
# Split the bottom pane (1) horizontally, creating pane 3 to its right
tmux split-window -h -t $SESSION_NAME:Docker.1

# Now we have 4 panes, indexed 0, 2, 1, 3 (top-left, top-min, top-right, bottom)
# Send the commands to the correct panes
tmux send-keys -t $SESSION_NAME:Docker.0 "docker logs -f bookapp-backend" C-m
tmux send-keys -t $SESSION_NAME:Docker.2 "docker logs -f bookapp-frontend" C-m
tmux send-keys -t $SESSION_NAME:Docker.1 "docker logs -f bookapp-proxy" C-m
tmux send-keys -t $SESSION_NAME:Docker.3 "docker ps -a" C-m

# Attach to the session to view it
tmux attach-session -t $SESSION_NAME