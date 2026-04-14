import json
from os import listdir
from os.path import join

file_names = [f for f in listdir("./recipes-json") if f.endswith(".json")]

slugs = {}

for file_name in file_names:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

    recipe["title"] = recipe["title"].lower()
    recipe["vod_title"] = recipe["vod_title"].lower()

    with open(file_path, "w", encoding="utf-8") as f:
        json.dump(recipe, f, ensure_ascii=False, indent=4)
