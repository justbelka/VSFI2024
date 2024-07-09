import React, { useEffect, useState, useContext } from 'react';
import { Box, Typography, Card, CardMedia, CardContent, CardActions, Button, Grid } from '@mui/material';
import { AuthContext } from '../contexts/AuthContext';
import { toast } from 'react-toastify';

const PremiumShishaList = () => {
  const { user, updateUserBalance } = useContext(AuthContext);
  const [images, setImages] = useState([]);
  const [purchasedImageIDs, setPurchasedImageIDs] = useState([]);

  useEffect(() => {
    fetch('/api/prem-images')
      .then((response) => response.json())
      .then((data) => setImages(data))
      .catch((error) => console.error('Error fetching images:', error));

    if (user) {
      fetch(`/api/purchased/ids/${user.username}`)
        .then(response => response.json())
        .then(data => {
          if (Array.isArray(data)) {
              setPurchasedImageIDs(data);
          } else {
              console.error('Unexpected response format:', data);
          }
      })
      .catch(error => console.error('Error fetching purchased image IDs:', error));
    }
  }, [user]);

  const handlePurchase = (imageId) => {
    fetch('/api/purchase', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        image_id: imageId,
        user_name: user.username,
      }),
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        }
        throw new Error('Failed to purchase image');
      })
      .then((data) => {
        toast.success('Image purchased successfully!');
        updateUserBalance();
        setPurchasedImageIDs([...purchasedImageIDs, imageId]);
      })
      .catch((error) => {
        toast.error(error.message);
        console.error('Error purchasing image:', error);
      });
  };

  return (
    <Box sx={{ paddingTop: 2 }}>
      <Grid container spacing={2}>
        {images.map((image) => (
          <Grid item key={image.id} xs={12} sm={6} md={4}>
            <Card style={{backgroundColor: "#343917"}}>
              <CardMedia
                component="img"
                height="200"
                image={image.url}
                alt={image.name}
              />
              <CardContent>
                <Typography variant="h6">{image.name}</Typography>
                <Typography variant="body2" color="textSecondary">
                  Price: {image.price}
                </Typography> <Typography variant="body2" color="textSecondary">
                  ID: {image.id}
                </Typography>
              </CardContent>
              <CardActions>
                {user &&  !purchasedImageIDs.includes(image.id.toString()) && (
                  <Button size="small" color="primary" onClick={() => handlePurchase(image.id)}>
                    Buy for 25 coins
                  </Button>
                )}
                {purchasedImageIDs.includes(image.id.toString()) && (
                  <Typography variant="body3" color="primary">
                    Purchased by you
                  </Typography>
                )}
              </CardActions>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default PremiumShishaList;
