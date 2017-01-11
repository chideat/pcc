#!/usr/bin/env bash

echo curl http://127.0.0.1:7030/api/v1/feeds/258/like?user_id=257 -d mood="love"
curl http://127.0.0.1:7030/api/v1/feeds/258/like?user_id=257 -d mood="love"

# echo curl -X DELETE http://127.0.0.1:7030/api/v1/feeds/258/like?user_id=257
# curl -X DELETE http://127.0.0.1:7030/api/v1/feeds/258/like?user_id=257
