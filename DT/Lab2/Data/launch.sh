sed 's/[\"\"\«\»\,\.\_\-\!\?\:-]\{1,\}/ /g' |
	sed 's/ \{1,\}/\n/g' |
       	sed '/^$/d'
