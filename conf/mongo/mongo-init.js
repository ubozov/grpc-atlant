db.createUser({
    user: '${DB_USER}',
    pwd: '${DB_PASSWORD}',
    roles: [
      {
        role: 'readWrite',
        db: '${DB_NAME}'
      }
    ]
});

db.products.createIndex({name:1}, {unique: true});