import os
import signal
import subprocess
import time

workDir = "."
sub = []

fileName = "dataSrv.exe"
apiName = "apiSrv.exe"
def main():
    print(os.getcwd())
    cmd = "go build -o {} ./dataServer/dataServer.go ".format(fileName)
    subprocess.Popen(cmd, shell=False)
    cmd = "go build -o {} ./apiServer/apiServer.go ".format(apiName)
    subprocess.Popen(cmd, shell=False)
    time.sleep(2)
    for i in range (6):
        cmd = "./{} --port 909{}".format(fileName,i)
        sub.append(subprocess.Popen(cmd, shell=False,stdout=subprocess.PIPE))
        print("data Server work on localhost:909{}".format(i))


    for i in range(2):
        cmd = "./{} --port 1908{}".format(apiName,i)
        sub.append(subprocess.Popen(cmd, shell=False,stdout=subprocess.PIPE))
        print("api Server work on localhost:1908{}".format(i))


def myHandler(signum, frame):
    print("receive ", signum)
    for p in sub:
        p.kill()
        print("kill")
    subprocess.Popen("del .\{}".format(fileName), shell=True)
    subprocess.Popen("del .\{}".format(apiName), shell=True)
    exit(0)


if __name__ == '__main__':
    signal.signal(signalnum=signal.SIGINT,handler=myHandler)
    main()
    while True:
        time.sleep(5)
