# Data Tester checks data in a database

## Install
```console
git clone https://github.com/mixo/data-tester.git
cd data-tester
docker build -t data-tester:latest .
```  

## Run
Copy the .env.example file:  
```console
cp .env.example .env
```

Change the env variable values.  

Run:  
```console
docker run --env-file=.env --rm data-tester:latest day-fluctuation -tn datatester_fixture -dc date -nc int_param,float_param -gc group_param -di -1 -md 40 -nd 10
```
