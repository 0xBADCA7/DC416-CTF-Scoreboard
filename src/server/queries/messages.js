const allMessagesQ = `
select id, posted, content
from messages;
`


const all = (db) => {
  return new Promise((resolve, reject) => {
    db.all(allMessagesQ, (err, rows) => {
      if (err) {
        reject(err)
      } else {
        resolve(rows)
      }
    })
  })
}


exports.all = all