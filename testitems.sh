tppd tx tpp create-item 'Rolex Submariner 1997 Gray' 'Rolex Submariner in good condition, it has no visible scratches and still works great. Bought in 1998, model year is 1997. It is the gray edition' 5 '40.741895,-73.989308' 3 watch,submariner,rolex 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 145 3 1 'Great Photos!' 0 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 145 3 1 'Great Photos!' 0 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 145 3 1 'Great Photos!' 0 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 0 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 450tpp
tppd tx tpp item-transferable 1 0 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-buyer 0 145 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transfer 0 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-rating 5 'Great seller, fast delivery' 0 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-resell 0 22 11 Nijmegen be,nl 'This watch did not fit me, so I resell this watch'  --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp


tppd tx tpp create-item G-Shock 'gshock in good condition, it has no visible scratches and still works great' 5 '40.7127281,-74.0060152' 3 grey 5 nl,be,us 2 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 11 2 1 'Great Photos!' 1 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 11 2 1 'Id Buy!' 1 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 11 2 1 'Great Photos!' 1 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 1 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transferable 1 1 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-buyer 1 11 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp delete-buyer 1 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-rating 3 'Item was not like the photos, but seller was nice nontheless' 1 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp


tppd tx tpp create-item G-Shock 'g-shock in great condition, it has no visible scratches and still works great' 5 '52.3727598,4.8936041'  3 grey 5 nl,be,us 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 12 3 1 'Great Photos!' 2 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 12 3 1 'Id Buy!' 2 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 12 3 1 'Great Photos!' 2 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 2 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transferable 1 2 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-buyer 2 12 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp


tppd tx tpp create-item Rolex 'rolex submariner in good condition, it has no visible scratches and still works great. Is really valueable over time, but I want to get into crypto instead.' 5 '40.741895,-73.989308' 3 watch,submariner,rolex 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 77 3 1 'Great Photos!' 3 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 77 3 1 'Great Photos!' 3 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 77 3 1 'Great Photos!' 3 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 3 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transferable 1 3 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp


tppd tx tpp create-item 'Air Jordans III' 'jordans in good condition, it has no visible scratches and walks nicely. I stood in the water with it but i cleaned the dirtyness' 3 '40.741895,-73.989308' 3 shoes,jordan,nike 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 45 3 1 'Great Photos!' 4 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 45 3 1 'Great Photos!' 4 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 45 3 1 'Great Photos!' 4 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 4 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transferable 1 4 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Air Jordans III' 'jordans in good condition, it has no visible scratches and walks nicely. I stood in the water with it but i cleaned the dirtyness' 3 '40.741895,-73.989308' 3 shoes,jordan,nike 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 45 3 1 'Great Photos!' 5 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 45 3 1 'Great Photos!' 5 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 45 3 1 'Great Photos!' 5 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 5 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transferable 1 5 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp


tppd tx tpp create-item 'Air Jordans III' 'jordans in good condition, it has no visible scratches and walks nicely. I stood in the water with it but i cleaned the dirtyness' 3 '40.741895,-73.989308' 3 shoes,jordan,nike 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 45 3 1 'Great Photos!' 6 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 45 3 1 'Great Photos!' 6 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 45 3 1 'Great Photos!' 6 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 6 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transferable 1 6 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-buyer 6 48 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-shipping 1 6 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-rating 5 'Great seller, fast delivery' 6 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-resell 6 23 22 '40.741895,-73.989308' us,be,nl 'This watch did not fit me, so I resell this watch'  --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item G-Shock 'gshock in good condition, it has no visible scratches and still works great' 5 '40.7127281,-74.0060152' 3 grey 5 nl,be,us 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 12 3 1 'Great Photos!' 7 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 12 3 1 'Id Buy!' 7 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 12 3 1 'Great Photos!' 7 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 7 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transferable 1 7 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-buyer 7 12 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp item-transfer 7 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-rating 3 'Item was not like the photos, but seller was nice nontheless' 7 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp


tppd tx tpp create-item G-Shock 'gshock in good condition, it has no visible scratches and still works great' 5 '40.7127281,-74.0060152'  3 grey 5 nl,be,us 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 14 3 1 'Great Photos!' 8 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 14 3 1 'Id Buy!' 8 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 14 3 1 'Great Photos!' 8 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 8 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transferable 1 8 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item G-Shock 'gshock in good condition, it has no visible scratches and still works great' 5 '40.7127281,-74.0060152'  3 grey 5 nl,be,us 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 15 3 1 'Great Photos!' 9 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 15 3 1 'Great Photos!' 9 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 15 3 1 'Great Photos!' 9 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 9 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transferable 1 9 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-buyer 9 15 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item Rolex 'submariner in good condition, it has no visible scratches and still works great ' 5 '40.741895,-73.989308' 3 watch,submariner,rolex 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 99 3 1 'Great Photos!' 10 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 99 3 1 'Great Photos!' 10 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 99 3 1 'Great Photos!' 10 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 10 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 10 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-buyer 10 99 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transfer 10 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item Rolex 'submariner in good condition, it has no visible scratches and still works great ' 5 '40.741895,-73.989308' 3 watch,submariner,rolex 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 112 3 1 'Great Photos!' 11 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 112 3 1 'Great Photos!' 11 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 112 3 1 'Great Photos!' 11 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 11 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 11 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
 
tppd tx tpp create-item 'Orient Bambino v4' 'Orient Bambino in decent condition, it has a few scratches and still works fine ' 5 '40.741895,-73.989308' 3 watch,orient,bambino 5 nl,be,de 2 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 8 2 1 'Great Photos!' 12 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 8 2 1 'Great Photos!' 12 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 8 2 1 'Great Photos!' 12 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 12 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 12 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-buyer 12 8 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transfer 12 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Orient Bambino v4' 'Orient Bambino in decent condition, it has a few scratches and still works fine ' 5 '40.741895,-73.989308' 3 watch,orient,bambino 5 nl,be,de 2 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 8 2 1 'Great Photos!' 13 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 8 2 1 'Great Photos!' 13 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 8 2 1 'Great Photos!' 13 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 13 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 13 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-buyer 13 8 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transfer 13 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Orient Bambino v4' 'Orient Bambino in decent condition, it has a few scratches and still works fine ' 5 '40.741895,-73.989308' 3 watch,orient,bambino 5 nl,be,de 2 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 9 2 1 'Great Photos!' 14 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 9 2 1 'Great Photos!' 14 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 9 2 1 'Great Photos!' 14 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 14 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 14 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-buyer 14 9 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp item-transfer 14 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Orient Bambino v4' 'Orient Bambino in decent condition, it has a few scratches and still works fine ' 5 '40.741895,-73.989308' 3 watch,orient,bambino 5 nl,be,de 2 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Rolex Submariner (no date)' 'Rolex submariner in good condition, it has a few scratches and but works great. Build year 2000' 5 '40.741895,-73.989308' 3 watch,submariner,rolex 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 102 3 1 'Great Photos!' 16 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 102 3 1 'Great Photos!' 16 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 102 3 1 'Great Photos!' 16 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 16 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 16 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Rolex Submariner (no date)' 'Rolex submariner in good condition, it has a few scratches but works great. Build year 2000' 5 '40.741895,-73.989308' 3 watch,submariner,rolex 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 112 3 1 'Great Photos!' 17 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 112 3 1 'Great Photos!' 17 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 112 3 1 'Great Photos!' 17 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 17 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 17 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Rolex Submariner' 'Rolex submariner in good condition, it has a few scratches but works great. Build year is 2000' 5 '40.741895,-73.989308' 3 watch,submariner,rolex 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 112 3 1 'Great Photos!' 18 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 112 3 1 'Great Photos!' 18 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 112 3 1 'Great Photos!' 18 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 18 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 18 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Omega Seamaster Professional' 'Omega Seamaster Professional in good condition, almost no visable scratches and still works great. Build year is 2001' 5 '40.741895,-73.989308' 3 watch,Omega,Seamaster 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 89 3 1 'Great Photos!' 19 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 89 3 1 'Great Photos!' 19 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 88 3 1 'Great Photos!' 19 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 19 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 19 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Omega Seamaster Professional' 'Omega Seamaster Professional in good condition, almost no visable scratches and still works great. Build year is 2001' 5 '40.741895,-73.989308' 3 watch,Omega,Seamaster 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 89 3 1 'Great Photos!' 20 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 89 3 1 'Great Photos!' 20 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 88 3 1 'Great Photos!' 20 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 20 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 20 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Rolex Daytona Cosmograph' 'Omega Seamaster Professional in good condition, almost some visable spots and scratches but works great.' 5 '40.741895,-73.989308' 3 watch,Rolex,Daytona,Cosmograph 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 345 3 1 'Great Photos!' 21 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 289 3 1 'Great Photos!' 21 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 288 3 1 'Great Photos!' 21 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 21 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 21 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp

tppd tx tpp create-item 'Rolex Daytona Cosmograph' 'Omega Seamaster Professional in good condition, almost some visable spots and scratches but works great.' 5 '40.741895,-73.989308' 3 watch,Rolex,Daytona,Cosmograph 5 nl 3 "" --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 345 3 1 'Great Photos!' 22 --from=user4 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 289 3 1 'Great Photos!' 22 --from=user2 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp create-estimation 288 3 1 'Great Photos!' 22 --from=user3 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp reveal-estimation 22 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp
tppd tx tpp  item-transferable 1 22 --from=user1 -y --chain-id=tpp --keyring-backend test --fees 150tpp


tppd q tpp list-item

