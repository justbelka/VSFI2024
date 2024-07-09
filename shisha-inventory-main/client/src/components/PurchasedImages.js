import React, { useEffect, useState, useContext } from 'react';
import { Box, Typography, Grid, Card, CardMedia, CardContent } from '@mui/material';
import { AuthContext } from '../contexts/AuthContext';

const PurchasedImages = () => {
    const { user } = useContext(AuthContext);
    const [images, setImages] = useState([]);

    useEffect(() => {
        if (user) {
            fetch(`/api/purchased/${user.username}`)
                .then((response) => response.json())
                .then((data) => setImages(data))
                .catch((error) => console.error('Error fetching purchased images:', error));
        }
    }, [user]);

    return (
        <Box sx={{ paddingTop: 2 }}>
            {images === null ? (
                <Typography variant="h6" color="textSecondary" align="center">
                    No purchased images found.
                </Typography>
            ) : (
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
                                        Buy time: {new Date(image.buytime).toLocaleString()}
                                    </Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                    ))}
                </Grid>
            )}
        </Box>
    );
};

export default PurchasedImages;
