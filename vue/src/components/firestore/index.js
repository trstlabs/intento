//import firebase from '../../main.js'
//import 'firebase/firestore'
//import 'firebase/database'
//import 'firebase/storage'

/*const config =
	process.env.NODE_ENV === 'development'
		? JSON.parse(process.env.VUE_APP_FIREBASE_CONFIG)
		: JSON.parse(process.env.VUE_APP_FIREBASE_CONFIG_PUBLIC)*/

//firebase.initializeApp(process.env.VUE_APP_FIREBASE_CONFIG_PUBLIC)
console.log("dsafsadf")
//export const firebase = firebase
export const db = firebase.firestore()
export const storageRef = firebase.storage().ref()
export const databaseRef = firebase.database()
export const usersRef = db.collection('users')
export const roomsRef = db.collection('chatRooms')
export const messagesRef = roomId => roomsRef.doc(roomId).collection('messages')

export const filesRef = storageRef.child('files')

export const dbTimestamp = firebase.firestore.FieldValue.serverTimestamp()
export const deleteDbField = firebase.firestore.FieldValue.delete()

