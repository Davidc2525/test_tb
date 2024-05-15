clone
	
 	git clone https://github.com/Davidc2525/test_tb 
	cd test_tb
use
  
    ./main help #show this help msg
    
    ./main new <account_number> #create new account
  
    ./main get <account_number>
  
    ./main create_bank <account_number>
  
    ./main transfer <debit_account_number> <credit_account_number> <amount> # use the bank account to transfer
  
    ./loop <debit_account_number>  <credit_account_number> <ammount> <iterate>  #iterate x 8000 

If it is the first time you must create the file from the tiger directory

	cd tiger

and create the file

	./tigerbeetle format --cluster=0 --replica=0 --replica-count=1 0_0.tigerbeetle

then in that same directory you execute TB

	./tigerbeetle start --addresses=3000 0_0.tigerbeetle
