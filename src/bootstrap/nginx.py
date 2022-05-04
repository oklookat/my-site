import os, utils, data, config

DEV_MIME_PATH = config.get("nginx.linux.dev_mime")
DEV_CONFIG_PATH = config.get("nginx.linux.dev_config")
if utils.getOS() == utils.OS_WINDOWS:
    DEV_MIME_PATH = config.get("nginx.windows.dev_mime")
    DEV_CONFIG_PATH = config.get("nginx.windows.dev_config")

def symlink_dev_config():
    utils.removeFile(DEV_CONFIG_PATH)
    os.symlink(utils.NGINX_DEV_CONF_PATH, DEV_CONFIG_PATH)
    utils.logger("nginx", "done")

def symlink_dev_mime():
    utils.removeFile(DEV_MIME_PATH)
    os.symlink(utils.NGINX_MIME_TYPES_PATH, DEV_MIME_PATH)
    utils.logger("nginx", "done")
