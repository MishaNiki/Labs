CREATE OR REPLACE FUNCTION "DTM"."AddStatistics"(text)
RETURNS void AS $$
DECLARE
str ALIAS FOR $1;
arr text[];
prow "DTM"."Lab2Stat"%ROWTYPE;
prow2 "DTM"."Lab2Stat"%ROWTYPE;
i integer;

BEGIN
	arr = string_to_array(str, ' ', '');
	
	i = 1;
	
	WHILE i <= array_length(arr, 1) LOOP

		SELECT INTO prow * FROM "DTM"."Lab2Stat"
			WHERE word = arr[i];
	
		IF prow IS NULL  THEN
			prow.word = arr[i];
			prow.ok = arr[i+1];
			prow.spam = arr[i+2];
 			INSERT INTO "DTM"."Lab2Stat"(
 	 			word, ok, spam) VALUES (prow.word, prow.ok + 1, prow.spam + 1);
		ELSE
			prow2.ok = arr[i+1];
			prow2.spam = arr[i+2];
			UPDATE "DTM"."Lab2Stat"
				SET ok = prow.ok + prow2.ok, spam = prow.spam + prow2.spam
				WHERE word = prow.word;
		END IF;
		i = i + 3;
	END LOOP;
END;
$$ LANGUAGE 'plpgsql';