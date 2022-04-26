import os, utils

def symlink_config():
    CONFIG_PATH = utils.get_absolute_by_relative("./data/nginx.conf")
    ORIGINAL_CONFIG_PATH = "/etc/nginx/nginx.conf"
    if utils.getOS() == "NT":
        ORIGINAL_CONFIG_PATH = "D:/Programming/nginx/conf/nginx.conf"
    utils.removeFile(ORIGINAL_CONFIG_PATH)
    os.symlink(CONFIG_PATH, ORIGINAL_CONFIG_PATH)
    
def symlink_mime():
    MIME_PATH = utils.get_absolute_by_relative("./data/mime.types")
    ORIGINAL_MIME_PATH = "/etc/nginx/mime.types"
    if utils.getOS() == "NT":
        ORIGINAL_MIME_PATH = "D:/Programming/nginx/conf/mime.types"
    utils.removeFile(ORIGINAL_MIME_PATH)
    os.symlink(MIME_PATH, ORIGINAL_MIME_PATH)