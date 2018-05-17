SERVER=$1
USER=$2
PASS=$3
# Error could be due to bitbucket locking your session, log in via ui and enter CAPTCHA
curl -u ${USER}:${PASS} -k -s https://${SERVER}/rest/api/1.0/projects/ | jq -r ".values[].key" | while read line
do
   curl -u ${USER}:${PASS} -k -s https://${SERVER}/rest/api/1.0/projects/${line}/repos | jq -r ".values[].links.clone[] | select(.name == \"http\") | .href" >> repository_list.txt
done

