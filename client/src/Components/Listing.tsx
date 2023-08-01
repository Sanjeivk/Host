import React, { useState, useEffect } from "react";

const people = [
  {
    name: "Lindsay Walton",
    role: "Front-end Developer",
    imageUrl:
      "https://images.unsplash.com/photo-1517841905240-472988babdf9?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=8&w=1024&h=1024&q=80",
    twitterUrl: "#",
    linkedinUrl: "#",
  },
  {
    name: "Lindsay Walton",
    role: "Front-end Developer",
    imageUrl:
      "https://images.unsplash.com/photo-1517841905240-472988babdf9?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=8&w=1024&h=1024&q=80",
    twitterUrl: "#",
    linkedinUrl: "#",
  },
  {
    name: "Lindsay Walton",
    role: "Front-end Developer",
    imageUrl:
      "https://images.unsplash.com/photo-1517841905240-472988babdf9?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=8&w=1024&h=1024&q=80",
    twitterUrl: "#",
    linkedinUrl: "#",
  },
  {
    name: "Lindsay Walton",
    role: "Front-end Developer",
    imageUrl:
      "https://images.unsplash.com/photo-1517841905240-472988babdf9?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=8&w=1024&h=1024&q=80",
    twitterUrl: "#",
    linkedinUrl: "#",
  },
  // More people...
];

export default function Listing() {
  const [listings, setListings] = useState([]);

  function formatDateString(dateString: string): string {
    const date = new Date(dateString);
    const options: Intl.DateTimeFormatOptions = {
      year: "numeric",
      month: "long",
      day: "numeric",
    };
    return date.toLocaleDateString(undefined, options);
  }

  useEffect(() => {
    // fetch("http://localhost:3000/listings")
    //   .then((response) => response.json())
    //   .then((data) => {
    //     setListings(data);
    //   });

    // declare the data fetching function
    const fetchData = async () => {
      const response = await fetch("http://localhost:3000/listings", {
        method: "GET",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      });
      //   setListings(response.json());
      const listings = await response.json();
      console.log(listings);
      setListings(listings);
    };

    // call the function
    fetchData()
      // make sure to catch any error
      .catch(console.error);
  }, []);

  return (
    <div className=" bg-white py-10">
      <div className="mx-auto max-w-7xl">
        <ul
          role="list"
          className="mx-auto grid max-w-2xl grid-cols-1 gap-x-6 gap-y-16 sm:grid-cols-2 lg:mx-0 lg:max-w-none lg:grid-cols-4"
        >
          {listings &&
            listings.map((listing: any) => (
              <li key={listing.id}>
                <img
                  className="aspect-[6/6] w-full rounded-2xl object-cover"
                  src={
                    "https://images.unsplash.com/photo-1517841905240-472988babdf9?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=8&w=1024&h=1024&q=80"
                  }
                  alt=""
                />
                <div className="justify between mt-3 flex w-full flex-row items-center justify-between">
                  <h3 className="text-lg font-semibold leading-8 tracking-tight text-gray-900">
                    {listing.occasion}
                  </h3>
                  <h3>{listing.review}</h3>
                </div>

                <h3 className="text-lg font-semibold leading-8 tracking-tight text-gray-900">
                  {listing.city} {listing.state} {listing.country}
                </h3>
                <p className="text-base leading-7 text-gray-600">
                  {formatDateString(listing.event_date)}
                </p>
              </li>
            ))}
        </ul>
      </div>
    </div>
  );
}
