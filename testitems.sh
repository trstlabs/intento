
trstd tx item create-item Rolex 'submariner in good condition, it has no visible scratches and still works great ' watch,submariner,rolex --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/rolex.jpeg?alt=media&token=b95b8a3c-1620-44fe-8eed-dfeebc57e1fc", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/rolex2.jpeg?alt=media&token=43ddf23c-b7db-42ad-9ae9-ef5c7ee0566b" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 14000000 5000000 1 'Great Photos!' 1 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 14000000 5000000 1 'Great Photos!' 1 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 14000000 5000000 1 'Great Photos!' 1 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 1 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transferable 1 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item G-Shock 'g-shock in great condition, it has no visible scratches and still works great' watch,gshock --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock.jpeg?alt=media&token=5b108efa-2ba1-429f-93d9-6be38dba9645", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock3.jpeg?alt=media&token=3e2f77ff-f084-4a7e-bf02-e84722b5e622", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock2.jpeg?alt=media&token=864c80c2-24e9-44d4-afd3-805a557836ee" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 12000000 5000000 1 'Great Photos!' 2 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 12000000 5000000 1 'Id Buy!' 2 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 12000000 5000000 1 'Great Photos!' 2 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 2 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transferable 2 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item buy-item 2 12000000 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item delete-prepayment 2 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-rating 2 'Item was not like the photos, but seller was nice nontheless' 4 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item Rolex 'rolex submariner in good condition, it has no visible scratches and still works great. Is really valueable over time, but I want to get into crypto instead.' watch,submariner,rolex --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock.jpeg?alt=media&token=5b108efa-2ba1-429f-93d9-6be38dba9645", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock3.jpeg?alt=media&token=3e2f77ff-f084-4a7e-bf02-e84722b5e622", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock2.jpeg?alt=media&token=864c80c2-24e9-44d4-afd3-805a557836ee" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 77 5000000 1 'Great Photos!' 3 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 77 5000000 1 'Great Photos!' 3 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 77 5000000 1 'Great Photos!' 3 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 3 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transferable 3 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst


trstd tx item create-item 'Air Jordans III' 'jordans in good condition, it has no visible scratches and walks nicely. I stood in the water with it but i cleaned the dirtyness' sneakers,jordan,nike  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/rolex3.jpeg?alt=media&token=792f12dc-26c2-4dfa-80f0-76d72fcf8b85", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/rolex2.jpeg?alt=media&token=43ddf23c-b7db-42ad-9ae9-ef5c7ee0566b" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 45 5000000 1 'Great Photos!' 4 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 45 5000000 1 'Great Photos!' 4 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 45 5000000 1 'Great Photos!' 4 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 4 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transferable 4 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Air Jordans III' 'jordans in good condition, it has no visible scratches and walks nicely. I stood in the water with it but i cleaned the dirtyness' sneakers,jordan,nike  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/jordans.jpeg?alt=media&token=92b09fcf-9fd8-4dcc-bd05-c3fb23a5e571", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/jordans2.jpeg?alt=media&token=4059d0d7-a544-445d-a351-96ef1f637bd9"  --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 45 5000000 1 'Great Photos!' 5 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 45 5000000 1 'Great Photos!' 5 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 45 5000000 1 'Great Photos!' 5 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 5 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transferable 5 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst


trstd tx item create-item 'Air Jordans III' 'jordans in good condition, it has no visible scratches and walks nicely. I stood in the water with it but i cleaned the dirtyness' sneakers,jordan,nike  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/jordans.jpeg?alt=media&token=92b09fcf-9fd8-4dcc-bd05-c3fb23a5e571", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/jordans2.jpeg?alt=media&token=4059d0d7-a544-445d-a351-96ef1f637bd9"  --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 45 5000000 1 'Great Photos!' 6 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 45 5000000 1 'Great Photos!' 6 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 45 5000000 1 'Great Photos!' 6 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 6 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transferable 6 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item buy-item 6 48 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-shipping 1 6 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-rating 6 'Great seller, fast delivery'  5 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-resell 6 23 22 '40.741895,-73.989308' us,be,nl 'This watch did not fit me, so I resell this watch'  --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item G-Shock 'gshock in good condition, it has no visible scratches and still works great' watch,gshock --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/jordans.jpeg?alt=media&token=92b09fcf-9fd8-4dcc-bd05-c3fb23a5e571", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/jordans2.jpeg?alt=media&token=4059d0d7-a544-445d-a351-96ef1f637bd9" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 120000000 5000000 1 'Great Photos!' 7 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 120000000 5000000 1 'Id Buy!' 7 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 120000000 5000000 1 'Great Photos!' 7 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 7 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transferable 7 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item buy-item 7 120000000 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item item-transfer 7 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-rating 7  'Item was  like the photos, and seller was nice ' 5 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst


trstd tx item create-item G-Shock 'gshock in good condition, it has no visible scratches and still works great' watch,gshock  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos  "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock.jpeg?alt=media&token=5b108efa-2ba1-429f-93d9-6be38dba9645", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock3.jpeg?alt=media&token=3e2f77ff-f084-4a7e-bf02-e84722b5e622", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock2.jpeg?alt=media&token=864c80c2-24e9-44d4-afd3-805a557836ee" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst 
trstd tx item create-estimation 140000000 5000000 1 'Great Photos!' 8 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 140000000 5000000 1 'Id Buy!' 8 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 140000000 5000000 1 'Great Photos!' 8 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 8 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transferable 8 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item G-Shock 'gshock in good condition, it has no visible scratches and still works great' watch,gshock  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://marketplace.trustlesshub.com/img/brand/icon.png" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst 
trstd tx item create-estimation 140000000 5000000 1 'Great Photos!' 9 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 140000000 5000000 1 'Great Photos!' 9 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 140000000 5000000 1 'Great Photos!' 9 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 9 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transferable 9 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item buy-item 9 140000000 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item Rolex 'submariner in good condition, it has no visible scratches and still works great ' watch,submariner,rolex  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock.jpeg?alt=media&token=5b108efa-2ba1-429f-93d9-6be38dba9645", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock3.jpeg?alt=media&token=3e2f77ff-f084-4a7e-bf02-e84722b5e622", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock2.jpeg?alt=media&token=864c80c2-24e9-44d4-afd3-805a557836ee" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 990000000 5000000 1 'Great Photos!' 10 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 990000000 5000000 1 'Great Photos!' 10 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 990000000 5000000 1 'Great Photos!' 10 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 10 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 10 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item buy-item 10 990000000 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transfer 10 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item Rolex 'submariner in good condition, it has no visible scratches and still works great ' watch,submariner,rolex  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock.jpeg?alt=media&token=5b108efa-2ba1-429f-93d9-6be38dba9645", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock3.jpeg?alt=media&token=3e2f77ff-f084-4a7e-bf02-e84722b5e622", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/gshock2.jpeg?alt=media&token=864c80c2-24e9-44d4-afd3-805a557836ee" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 112000000 5000000 1 'Great Photos!' 11 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 112000000 5000000 1 'Great Photos!' 11 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 112000000 5000000 1 'Great Photos!' 11 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 11 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 11 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
 
trstd tx item create-item 'Orient Bambino v4' 'Orient Bambino in decent condition, it has a few scratches and still works fine ' watch,Orient,Bambino   --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/rolex.jpeg?alt=media&token=b95b8a3c-1620-44fe-8eed-dfeebc57e1fc", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/rolex2.jpeg?alt=media&token=43ddf23c-b7db-42ad-9ae9-ef5c7ee0566b"  --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 80000000 5000000 1 'Great Photos!' 12 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 80000000 5000000 1 'Great Photos!' 12 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 80000000 5000000 1 'Great Photos!' 12 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 12 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 12 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item buy-item 12 80000000 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transfer 12 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Orient Bambino v4' 'Orient Bambino in decent condition, it has a few scratches and still works fine ' watch,Orient,Bambino  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/rolex2.jpeg?alt=media&token=43ddf23c-b7db-42ad-9ae9-ef5c7ee0566b", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/rolex.jpeg?alt=media&token=b95b8a3c-1620-44fe-8eed-dfeebc57e1fc" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 80000000 5000000 1 'Great Photos!' 13 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 80000000 5000000 1 'Great Photos!' 13 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 80000000 5000000 1 'Great Photos!' 13 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 13 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 13 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item buy-item 13 8 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transfer 13 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Orient Bambino v4' 'Orient Bambino in decent condition, it has a few scratches and still works fine ' watch,Orient,Bambino  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/bambino.jpeg?alt=media&token=ad64dc9c-151d-44be-be9a-6287b2166f8c", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/bambino2.jpeg?alt=media&token=ec599c09-c585-434d-b924-7ec38b9db9e9" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 9000000 5000000 1 'Great Photos!' 14 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 9000000 5000000 1 'Great Photos!' 14 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 9000000 5000000 1 'Great Photos!' 14 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 14 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 14 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item buy-item 14 9000000 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transfer 14 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Orient Bambino v4' 'Orient Bambino in decent condition, it has a few scratches and still works fine ' watch,Orient,Bambino  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/bambino3.jpeg?alt=media&token=42ddcd47-2d77-4f60-909a-8eedf559a912", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/bambino.jpeg?alt=media&token=ad64dc9c-151d-44be-be9a-6287b2166f8c", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/bambino4.jpeg?alt=media&token=e0b1752a-1a81-4e64-b427-1afd8c7c4625"  --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Rolex Submariner (no date)' 'Rolex submariner in good condition, it has a few scratches and but works great. Build year 2000' watch,submariner,rolex   --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/rolex.jpeg?alt=media&token=b95b8a3c-1620-44fe-8eed-dfeebc57e1fc", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/bambino2.jpeg?alt=media&token=ec599c09-c585-434d-b924-7ec38b9db9e9"  --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 1020000000 5000000 1 'Great Photos!' 16 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 1020000000 5000000 1 'Great Photos!' 16 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 1020000000 5000000 1 'Great Photos!' 16 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 16 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 16 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Rolex Submariner (no date)' 'Rolex submariner in good condition, it has a few scratches but works great. Build year 2000' watch,submariner,rolex  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos  "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub.jpeg?alt=media&token=5f1704e3-cf4f-4b9e-b390-caeb2642af48", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub2.jpeg?alt=media&token=d2f4eb71-7a7b-42cf-b738-3059748f76bf", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub3.jpeg?alt=media&token=e7213554-0b0e-4f6d-be08-c03c2a8019c5", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub4.jpeg?alt=media&token=23aefbfd-7cb3-44e9-8dd5-8861c78e43f4" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 1120000000 5000000 1 'Great Photos!' 17 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 1120000000 5000000 1 'Great Photos!' 17 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 1120000000 5000000 1 'Great Photos!' 17 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 17 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 17 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Rolex Submariner' 'Rolex submariner in good condition, it has a few scratches but works great. Build year is 2000' watch,submariner,rolex  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos  "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub.jpeg?alt=media&token=5f1704e3-cf4f-4b9e-b390-caeb2642af48", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub2.jpeg?alt=media&token=d2f4eb71-7a7b-42cf-b738-3059748f76bf", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub3.jpeg?alt=media&token=e7213554-0b0e-4f6d-be08-c03c2a8019c5", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub4.jpeg?alt=media&token=23aefbfd-7cb3-44e9-8dd5-8861c78e43f4"  --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 1120000000 5000000 1 'Great Photos!' 18 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 1120000000 5000000 1 'Great Photos!' 18 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 1120000000 5000000 1 'Great Photos!' 18 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 18 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 18 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Omega Seamaster Professional' 'Omega Seamaster Professional in good condition, almost no visable scratches and still works great. Build year is 2001' watch,Omega,seamaster  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos  "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub3.jpeg?alt=media&token=e7213554-0b0e-4f6d-be08-c03c2a8019c5", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub4.jpeg?alt=media&token=23aefbfd-7cb3-44e9-8dd5-8861c78e43f4", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub.jpeg?alt=media&token=5f1704e3-cf4f-4b9e-b390-caeb2642af48", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/sub2.jpeg?alt=media&token=d2f4eb71-7a7b-42cf-b738-3059748f76bf"  --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 8900000000 5000000 1 'Great Photos!' 19 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 8900000000 5000000 1 'Great Photos!' 19 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 8900000000 5000000 1 'Great Photos!' 19 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 19 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 19 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Omega Seamaster Professional' 'Omega Seamaster Professional in good condition, almost no visable scratches and still works great. Build year is 2001' watch,Omega,seamaster   --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/omega.jpeg?alt=media&token=2c623a1b-fbf3-4976-a783-c598c63afa6f", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/omega2.jpeg?alt=media&token=cb0a64d8-a934-4272-a033-46897e74ee3b" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 8900000000 5000000 1 'Great Photos!' 20 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 8900000000 5000000 1 'Great Photos!' 20 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 8900000000 5000000 1 'Great Photos!' 20 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 20 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 20 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Rolex Daytona Cosmograph' 'Rolex daytona Professional in good condition, almost some visable spots and scratches but works great.' watch,Daytona,Cosmograph,rolex  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/omega.jpeg?alt=media&token=2c623a1b-fbf3-4976-a783-c598c63afa6f", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/omega2.jpeg?alt=media&token=cb0a64d8-a934-4272-a033-46897e74ee3b" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 34500000000 5000000 1 'Great Photos!' 21 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 28900000000 5000000 1 'Great Photos!' 21 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 28800000000 5000000 1 'Great Photos!' 21 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 21 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item  item-transferable 21 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst

trstd tx item create-item 'Rolex Daytona Cosmograph' 'Omega Seamaster Professional in good condition, almost some visable spots and scratches but works great.' watch,submariner,rolex  --deposit_amount 5000000  --condition 5 --local_pickup '40.741895,-73.989308'  --estimation_count 3 --shipping_region uk,be,nl --shipping_cost 5000000 --photos "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/daytona.jpeg?alt=media&token=a3aab3f0-2f35-4dd5-a416-c4b1503cb72b", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/daytona2.jpeg?alt=media&token=f5d5cf64-64df-4336-aa81-6b39c550c1e4", "https://firebasestorage.googleapis.com/v0/b/trustitems-cbb92.appspot.com/o/daytona3.jpeg?alt=media&token=5634c3c4-00d9-4ebe-984f-c0155ed9722d" --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 34500000000 5000000 1 'Great Photos!' 22 --from user4 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 28900000000 5000000 1 'Great Photos!' 22 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item create-estimation 24900000000 5000000 1 'Great Photos!' 22 --from user3 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item reveal-estimation 22 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst
trstd tx item item-transferable 22 --from user1 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst


trstd tx item create-item 'Donald Trump is greaaat' 'As you can see, the donald is back and looking greater than eva. Give me a diet coke please, thank you.' nft,donald,trump  --deposit_amount 5000000  --estimation_count 3 --token_uri https://meta.creepz.co/shapeshifters/2 --from user2 -y --chain-id trst_chain_1 --keyring-backend test --fees 150utrst


trstd q item list-items

