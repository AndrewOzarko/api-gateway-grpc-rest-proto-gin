db = db.getSiblingDB('gateway');
db.createUser({
  user: 'gateway',
  pwd: 'secret',
  roles: [{
    role: 'readWrite',
    db: 'gateway'
  }]
})