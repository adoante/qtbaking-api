import json
from os import listdir
from os.path import join

file_names = [f for f in listdir("./recipes-json") if f.endswith(".json")]

for file_name in file_names:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

        tags = recipe.get("tags")

        print(recipe.get("title"))
        print(tags)
        print("=================================================")

        tagsToRemove = []

        for tag in tags:
            ans = input(f"Remove '{tag}'? (y/n): ")

            if ans == "y":
                print(f"Removing '{tag}'")
                tagsToRemove.append(tag)
            if ans == "n":
                print(f"Keeping: '{tag}'")

        for tag in tagsToRemove:
            tags.remove(tag)
        print(tags)
        print("=================================================")

        while (True):
            ans = input("Add tags? (y/n): ")

            if ans == "y":
                tag = input("Enter tag: ")
                tags.append(tag)
                print(f"Adding '{tag}'")

            if ans == "n":
                break

            print(tags)

        print("****************************************************************************************")

    with open(file_path, "w", encoding="utf-8") as f:
        json.dump(recipe, f, ensure_ascii=False, indent=4)
