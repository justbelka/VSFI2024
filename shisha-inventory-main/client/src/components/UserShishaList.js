import React, { useEffect, useState } from 'react';
import { Box, Typography, Card, CardMedia, CardContent, Grid } from '@mui/material';

const UserShishaList = () => {
    const [images, setImages] = useState([]);

    useEffect(() => {
        fetch('/api/user-images')
            .then((response) => response.json())
            .then((data) => setImages(data))
            .catch((error) => console.error('Error fetching images:', error));
    }, []);

    return (
        <Box sx={{ paddingTop: 2 }}>
             {images === null ? (
                <Typography variant="h6" color="textSecondary" align="center">
                    No images found.
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
                                    Uploaded by: {image.owner}
                                </Typography>
                                <Typography variant="body2" color="textSecondary">
                                    Uploaded at: {new Date(image.uploadedAt).toLocaleString()}
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

export default UserShishaList;
