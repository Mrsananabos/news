SELECT n.id, n.title, n.content, ARRAY_AGG(nc.category_id) AS categories
FROM news n
         LEFT JOIN news_categories nc ON n.id = nc.news_id
GROUP BY n.id
    LIMIT $1 OFFSET $2;