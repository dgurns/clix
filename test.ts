interface Person {
	firstName: string;
	lastName: string;
}

function sayHello(p: Person): string {
	return `Hello, ${p.firstName}!`;
}

const me: Person = {
	firstName: 'Dan',
	lastName: 'Gurney',
};

console.log(sayHello(me));
