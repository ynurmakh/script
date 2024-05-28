# !bin/bash

for i in {1..100};do
    LOCK_STATUS=$(pgrep i3lock)
    
    if [ -n "$LOCK_STATUS" ]; then
        echo "Компьютер заблокирован"
    else
        echo "Компьютер не заблокирован"
    fi
    
    sleep 1
done