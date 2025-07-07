#!/bin/bash

jet \
  -source=postgres \
  -dsn="postgresql://api_database_jqur_user:jx7YxQwWkqOdLAavkYl53O3AFCSGRQzu@dpg-d1in3u0dl3ps73d35oqg-a.oregon-postgres.render.com:5432/api_database_jqur?sslmode=require" \
  -port=5432 \
  -user=api_database_jqur_user \
  -password=jx7YxQwWkqOdLAavkYl53O3AFCSGRQzu \
  -dbname=api_database_jqur \
   -schema=public \
  -path=./.gen
