SELECT A.* FROM test A,
(SELECT aid, MAX(day) max_day FROM test WHERE aid in (2,3) GROUP BY aid) B
WHERE A.aid = B.aid AND A.day = B.max_day
ORDER BY A.install DESC;
