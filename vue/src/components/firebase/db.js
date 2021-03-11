import firebase from "firebase";
import "firebase/firestore";
import 'firebase/storage'


var firebaseConfig = {
    apiKey: process.env.VUE_APP_API_KEY,
    authDomain: "trustitems-cbb92.firebaseapp.com",
    databaseURL: "https://trustitems-cbb92.firebaseio.com",
    projectId: "trustitems-cbb92",
    storageBucket: "trustitems-cbb92.appspot.com",
    messagingSenderId: "1033037336313",
    appId: "1:1033037336313:web:475aaf4e502bc36adf9c28",
    measurementId: "G-XQ9H3T2QNM",
  };





firebase.initializeApp(firebaseConfig);

export const db = firebase.firestore()
export const storageRef = firebase.storage().ref()
export const databaseRef = firebase.database()
export const usersRef = db.collection('users')
export const roomsRef = db.collection('chatRooms')
export const messagesRef = roomId => roomsRef.doc(roomId).collection('messages')

export const filesRef = storageRef.child('files')

export const dbTimestamp = firebase.firestore.FieldValue.serverTimestamp()
export const deleteDbField = firebase.firestore.FieldValue.delete()