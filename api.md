# API Endpoints

## Vods
- /vods/
- /vods?sort=created_at&order=asc|desc
 - order: desc by default
- /vods/:slug

## Recipes
- /recipes/
- /recipes?title=value&tag=value&match=exact|partial&limit=value&offset=value
 - match: partial by default
 - limit: 10 by default
 - offset: 0 by default
- /recipes/:id

## Components
- /components/

## Ingredients
- /ingredients/

## Tools
- /tools/

## Notes
- /notes/

## Tags
- /tags/

## Bakealongs
- /bakealongs/
- /bakealongs?sort=created_at&order=asc|desc&title=value&tag=value&match=exact|partial&limit=value&offset=value
 - match: partial by default
 - limit: 10 by default
 - offset: 0 by default
 - order: desc by default
- /bakealongs/:ytid 
