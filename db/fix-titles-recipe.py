import json
from os import listdir
from os.path import join

file_names = [f for f in listdir("./recipes-json") if f.endswith(".json")]

for file_name in file_names:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

        recipe["title"] = file_name.split(".")[0]

    with open(file_path, "w", encoding="utf-8") as f:
        json.dump(recipe, f, ensure_ascii=False, indent=4)
