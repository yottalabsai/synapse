#! /bin/bash
### BEGIN INIT INFO
# Provides:             period-jurassic
# Required-Start:       $syslog $remote_fs
# Required-Stop:        $syslog $remote_fs
# Should-Start:         $local_fs
# Should-Stop:          $local_fs
# Default-Start:        3 4 5
# Default-Stop:         0 1 2 6
# Short-Description:    Base Service
# Description:          The base service provided for the kafka, sms and email functions to the rest modules
### END INIT INFO

# Using the LSB functions to perform the operations
# . /lib/lsb/init-functions

# start the application
PIDFILE="./pid"
NAME=synapse
PROFILE=dev

### pull & compile the module
do_update()
{
        export PATH=$PATH:/usr/local/go/bin
        echo "[Start to pull the code from dev branch]"
        git checkout dev && git pull
        echo "[start to compile and build the application - $NAME]"
        # go mod tidy
        go mod download
        go mod vendor
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $NAME main.go
}

do_start()
{
        sleep 1
        echo "[Starting $NAME]......"
        nohup ./$NAME>/dev/null 2>&1 & echo $!>$PIDFILE
        echo "[$NAME started]"
}

do_console()
{
        sleep 1
        echo "[Starting $NAME]......"
        ./$NAME
}

do_stop()
{

        echo "[Stopping $NAME.......]"
        if [ -e $PIDFILE ]; then
                kill -9 $(cat $PIDFILE)
                rm -rf $PIDFILE
                echo "[$NAME stopped]"
        else
                echo "[Couldn't find service $NAME]"
        fi
}

case "$1" in
        console)
                if [ -e $PIDFILE ]; then
                        echo "[$NAME started already]"
                        exit
                fi
                #Start from console
                do_console
                ;;
        start)
                if [ -e $PIDFILE ]; then
                        echo "[$NAME started already]"
                        exit
                fi
                #Start the daemon
                do_start
                ;;
        stop)
                # Stop the daemon
                do_stop
                ;;
        restart)
                do_stop && sleep 5 && do_start
                ;;
        update)
                do_update
                rc=$?
                if [[ $rc -ne 0 ]] ; then
                        echo "[Build failed!]"; exit $rc
                fi
                do_stop && sleep 5 && do_start
                ;;
        status)
                if [ -e $PIDFILE ]; then
                        echo "[$NAME is running. PID is: $(cat $PIDFILE)]"
                else
                        echo "[$NAME is not running]"
                fi
                ;;
        *)
        echo $"Usage: $0 {console|start|stop|update|restart|status}"
        exit 1
esac