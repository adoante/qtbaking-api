import json
from os import listdir
from os.path import join

file_names = [f for f in listdir("./recipes-json") if f.endswith(".json")]

slugs = {}

for file_name in file_names:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

    video_url = recipe.get("video_url")

    if video_url:
        print(recipe["title"])
        print(video_url)

        recipe["slug"] = recipe["video_url"].split("/")[-1]

        if recipe["slug"] in slugs:
            slugs[recipe["slug"]]["count"] = slugs[recipe["slug"]]["count"] + 1
            slugs[recipe["slug"]]["recipes"].append(recipe["title"])
        else:
            slugs[recipe["slug"]] = {
                "count": 1,
                "recipes": [recipe["title"]]
            }

        print(recipe["slug"])

    else:
        recipe["slug"] = ""
        print("No video ID")

    with open(file_path, "w", encoding="utf-8") as f:
        json.dump(recipe, f, ensure_ascii=False, indent=4)

print("")
dups = ([slug for slug in slugs.keys() if slugs[slug]["count"] > 1])
for dup in dups:
    print(dup, slugs[dup])
