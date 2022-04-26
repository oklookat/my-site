import os, errno
import sys
import subprocess
import shutil
import ctypes


def throwFatalErr(err: Exception):
    print(""" 
        -------- ERROR --------
        {err}
        -------- ERROR --------
        """.format(err=err))
    sys.exit(1)


def moveFile(fromPath: str, toPath: str):
    try:
        os.rename(fromPath, toPath)
    except Exception as e:
        throwFatalErr(e)


def copyFile(fromPath: str, toPath: str):
    try:
        shutil.copy(fromPath, toPath, follow_symlinks=True)
    except Exception as e:
        throwFatalErr(e)


def removeDir(path: str):
    try:
        shutil.rmtree(path)
    except Exception as e:
        throwFatalErr(e)

def removeFile(path: str):
    try:
        os.remove(path)
    except OSError as e:
        if e.errno != errno.ENOENT:
            throwFatalErr(e)

def getOS():
    return os.name.upper()

def runCommand(command: str, exit_if_error: bool = True) -> str:
    try:
        out = subprocess.run(command, shell=True, check=True, capture_output=True)
        res = out.stdout.decode("UTF-8")
        return res
    except Exception as e:
        if exit_if_error == False:
            raise e
        throwFatalErr(e)


def run_mkcert_command(command: str) -> str:
    mkcertPath = "mkcert "
    if getOS() == "NT":
        mkcertPath = ".\data\mkcert "
    try:
        out = runCommand(mkcertPath + command, True)
        return out.strip()
    except Exception as e:
        print("failed to run mkcert, install it")


# get main.py dir
def getExecutionDir() -> str:
    cwd = os.getcwd() + "/" + "main.py"
    return cwd


def isAdmin() -> bool:
    is_admin = False
    try:
        is_admin = os.getuid() == 0
    except AttributeError:
        is_admin = ctypes.windll.shell32.IsUserAnAdmin() != 0
    return is_admin

def get_absolute_by_relative(path: str) -> str:
    return os.path.abspath(path)


class CommandsStack:
    def __init__(self):
        self.cmdList: list[str] = []

    def addCommand(self, command: str):
        self.cmdList.append(command)

    def runCommands(self):
        for command in self.cmdList:
            runCommand(command)
