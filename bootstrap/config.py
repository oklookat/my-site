import yaml, dpath.util
import utils

def get(key: str):
    parsed = None
    with open(utils.CONFIG_PATH) as config:
        parsed = yaml.safe_load(config)
    return dpath.util.get(parsed, key, ".")
