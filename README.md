# Basic-Hash-Cracking
Prime Suspects: Outpost24 Hash Cracking CTF WriteUp
Challenge Description:
Json's notorious ransomware gang has unveiled what they claim to be an uncrackable hashing algorithm, derived from the digits of the largest known prime number. Mocking the world, they have publicly shared the hash of their master password, confident no one can break the encryption.

Crypto analysts believe the password is composed of digits only, sourced from this massive prime. For maximum security, the gang ensured the password is at least 10 digits long, hidden deep within the prime’s vast expanse. If you manage to crack the hash, you also found the flag, wrap it in the format O24{cracked_hash_here}. The answer lies in the prime.

Hash: $2y$05$tJ5qkcBGrjiRfZZAlkSsP.kcVStH7oCzsery3nN1sgXk02xThNck6
Author: Mikael Svall

Tools Used
- Kali Linux
- Windows CMD
- Go Lang
- John The Ripper

Solution Approach:
The first step to find out what hashing algorithm was used. To achieve this there are a number of tools that can be used. For this challenge I used hashid and nth. 

Here are the commands;

hashid ‘given_hash’

nth -t given_hash

In this context that’ll be;

hashid ‘$2y$05$tJ5qkcBGrjiRfZZAlkSsP.kcVStH7oCzsery3nN1sgXk02xThNck6’
![Screenshot 2024-12-31 115743](https://github.com/user-attachments/assets/3471915c-fe6b-42e5-afbe-27cb6544b89c)

nth -t ‘$2y$05$tJ5qkcBGrjiRfZZAlkSsP.kcVStH7oCzsery3nN1sgXk02xThNck6’
![Screenshot 2024-12-31 120229](https://github.com/user-attachments/assets/f1ef58d5-58ba-4883-8c24-a538e6c80914)

Using nth helped confirm the exact hashing algorithm used. Now that we are certain that the hashing algorithm used is bcrypt. We will proceed to crack the hash. 

In order to crack the hash, we need to find a possible 10 digit sequence with a bcrypt that matches the given hash ($2y$05$tJ5qkcBGrjiRfZZAlkSsP.kcVStH7oCzsery3nN1sgXk02xThNck6). To make it easier, I wrote a Go program to extract all 10 digit sequences from the prime number file into a new file.

//extract_sequences.go, see repo for code//

Following this I wrote another Go program that will compare the hash of each sequence to the given hash. 

//hash_cracker.go, see repo for code//

Results and Conclusion
It took 2 hours 55 minutes and about 27 seconds to loop through 14.69 million sequences in order to find a successful match.
Flag: O24{9999903898}

![Screenshot 2024-12-31 111607](https://github.com/user-attachments/assets/0aee06d1-d6cb-4cc9-8e43-6ee98efb6f39)


It is important to note that the speed of cracking is heavily dependent on the hardware resources present on the host PC, which is why I executed this program on my host PC and not on my VM. 

An alternative to this would be to use JohnTheRipper on Kali with the following command; 
john --format=bcrypt hash.txt --wordlist=sequences.txt
