#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>
#include <errno.h>

#include <pthread.h>
#include <semaphore.h>

//#define CHAIRS 3
//#define CUSTOMERS 5

/* Initialize variables--counts, locks, and semaphores */
int waiting = 0;
int total = 0;
int numCustomers, numChairs;
sem_t mutex;
sem_t customer;
sem_t barber;
sem_t customerDone;
sem_t barberDone;


void cut_hair() {
	printf("The barber is cutting hair...\n");
	sleep(4);
}

void get_hair_cut(int num) {
	printf("Customer %d is getting their hair cut...\n", num);
}

void balk() {
	printf("The customer balks and leaves.\n");
}

// When a customer arrives, the customer procedure is executed which begins by acquiring mutex to enter a critical region. 
// Subsequent arriving customers have to wait until the first customer has released mutex.

#define customer_print(msg) printf("Customer %d (waiting: %d) -- " msg "\n", num, waiting)
#define barber_print(msg)   printf("Barber     (waiting: %d) -- " msg "\n", waiting)

void *customer_thread(void *id) {
	int num = *(int *)id + 1;
	
	// Acquire mutex to enter a critical region, protecting number of available seats
	sem_wait(&mutex);
	// After acquiring mutex, a customer checks to see if the number of waiting customers is less than the number of chairs.
	if (waiting == numChairs) {
  		// mutex lock on seats is released and the customer leaves without a haircut.
		customer_print("Seats are full.");
    	sem_post(&mutex);
    	balk();
    	// the number of customers that came by so far
    	total++;
    	pthread_exit(NULL);
 	} else {
		customer_print("Found open seat.");
		// If there is an available chair, the waiting counter is incremented
		waiting++;
		// the number of customers that came by so far
		total++;
	}

	// customer releases mutex
	sem_post(&mutex);

	// wake up barber if they're asleep
	sem_post(&customer);
	
	// wait for barber to show them to chair
	customer_print("Waiting for barber.");
	sem_wait(&barber);

	// Acquire mutex
	sem_wait(&mutex);
	// one less customer in waiting room
	waiting--;
	// releases mutex on seats
	sem_post(&mutex);

	// and begins the haircut
	get_hair_cut(num);

	// customer signals customerDone
	sem_post(&customerDone);
	// wait until barber says done
	sem_wait(&barberDone);

	pthread_exit(NULL);

}

// When the barber arrives for work, the barber procedure is executed blocking the barber on the customer semaphore until a customer arrives.
void *barber_thread(void *args) {

	while (!(total == numCustomers && waiting == 0)) {
		// wait for or get next customer       
		// customer signals
		barber_print("Ready for customer.");
		sem_wait(&customer);  // go to sleep if no customers

		// barber wakes, signals barber
		sem_post(&barber);

		// cut hair
		cut_hair();

		// tell customer that haircut is done
		sem_wait(&customerDone);
		sem_post(&barberDone);
		barber_print("Done cutting hair.");
	}
	printf("The barber has finished for the day.\n");
	printf("\nBarbershop Problem Completed.\n");
	pthread_exit(NULL);

}


int main(int argc, char *args[]) {
	int cu, ch;

	printf("Welcome to the barbershop problem!\n\n");
	printf("Enter number of customers: ");
	scanf("%d",&cu);
    numCustomers = cu;
    printf("Enter number of chairs: ");
    scanf("%d",&ch);
    numChairs = ch;
	
	int i, id[numCustomers];
	pthread_t b_thread;
    pthread_t c_threads[numCustomers];

    printf("\nInitializing semaphores...");
    sleep(1);
   	// Initialize lock
	sem_init(&mutex, 0, 1);
	// Create customer and barber
	sem_init(&customer, 0, 0);
	sem_init(&barber, 0, 0);
	// done semaphores
	sem_init(&customerDone, 0, 0);
	sem_init(&barberDone, 0, 0);
	printf("Done.\n\n");

	printf("Beginning Simulation.\n\n");
	srand((unsigned int) time(NULL));

	// Create the barber. 
    pthread_create(&b_thread, NULL, barber_thread, NULL); 

    // Create the customers. 
    for (i = 0; i < numCustomers; i++) {
    	id[i] = i;
        pthread_create(&c_threads[i], NULL, customer_thread, (void*)&id[i]);
        pthread_detach(c_threads[i]);  // Thread shouldn't hold resources
        sleep(rand() % 2 + 1);
    }
    pthread_exit(NULL);

	// join threads to wait for completion
	for (i = 0; i < numCustomers; i++) {
		pthread_join(c_threads[i], NULL);
		sleep(1);
	}

    // join barber thread to finish (blocks until thread exits)
    sem_post(&customerDone);
    pthread_join(b_thread, NULL);
	
	return 0;
}
