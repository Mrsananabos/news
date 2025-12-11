SELECT n.id,
       n.title,
       n.content,
       COALESCE(ARRAY_AGG(nc.category_id) FILTER (WHERE nc.category_id IS NOT NULL), '{}') AS categories
FROM news n
         LEFT JOIN news_categories nc ON n.id = nc.news_id
GROUP BY n.id
ORDER BY n.id DESC
    LIMIT $1 OFFSET $2;