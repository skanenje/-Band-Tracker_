import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { MapPin, Calendar, Users, Disc, Loader2 } from 'lucide-react';
import { Artist, Location, Date, Relation } from '../types';
import toast from 'react-hot-toast';

const ArtistDetails = () => {
  const { id } = useParams();
  const [artist, setArtist] = useState<Artist | null>(null);
  const [locations, setLocations] = useState<Location[]>([]);
  const [dates, setDates] = useState<Date[]>([]);
  const [relations, setRelations] = useState<Relation[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'info' | 'locations' | 'dates'>('info');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [artistRes, locationsRes, datesRes, relationsRes] = await Promise.all([
          fetch(`http://localhost:8080/artist/${id}`),
          fetch('http://localhost:8080/locations'),
          fetch('http://localhost:8080/dates'),
          fetch('http://localhost:8080/relations')
        ]);

        const [artistData, locationsData, datesData, relationsData] = await Promise.all([
          artistRes.json(),
          locationsRes.json(),
          datesRes.json(),
          relationsRes.json()
        ]);

        setArtist(artistData);
        setLocations(locationsData.index);
        setDates(datesData.index);
        setRelations(relationsData.index);
      } catch (error) {
        toast.error('Failed to fetch artist details');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  if (loading || !artist) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <Loader2 className="w-8 h-8 animate-spin text-indigo-500" />
      </div>
    );
  }

  const artistLocations = locations.find(loc => loc.id === artist.id)?.locations || [];
  const artistDates = dates.find(date => date.id === artist.id)?.dates || [];
  const artistRelations = relations.find(rel => rel.id === artist.id)?.datesLocations || {};

  return (
    <div className="max-w-4xl mx-auto">
      <div className="grid md:grid-cols-2 gap-8">
        <div className="relative group">
          <img
            src={artist.image}
            alt={artist.name}
            className="w-full rounded-lg shadow-xl"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent rounded-lg opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>

        <div>
          <h1 className="text-4xl font-bold mb-4">{artist.name}</h1>
          
          <div className="space-y-4">
            <div className="flex items-center space-x-2">
              <Users className="text-indigo-500" />
              <span className="text-gray-300">Members:</span>
              <div className="flex flex-wrap gap-2">
                {artist.members.map((member, index) => (
                  <span
                    key={index}
                    className="text-sm px-2 py-1 bg-gray-800 rounded-full"
                  >
                    {member}
                  </span>
                ))}
              </div>
            </div>

            <div className="flex items-center space-x-2">
              <Calendar className="text-indigo-500" />
              <span className="text-gray-300">Created:</span>
              <span>{artist.creationDate}</span>
            </div>

            <div className="flex items-center space-x-2">
              <Disc className="text-indigo-500" />
              <span className="text-gray-300">First Album:</span>
              <span>{artist.firstAlbum}</span>
            </div>
          </div>
        </div>
      </div>

      <div className="mt-8">
        <div className="flex space-x-4 mb-6">
          <button
            onClick={() => setActiveTab('info')}
            className={`px-4 py-2 rounded-lg transition-colors ${
              activeTab === 'info'
                ? 'bg-indigo-500 text-white'
                : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
            }`}
          >
            Info
          </button>
          <button
            onClick={() => setActiveTab('locations')}
            className={`px-4 py-2 rounded-lg transition-colors ${
              activeTab === 'locations'
                ? 'bg-indigo-500 text-white'
                : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
            }`}
          >
            Locations
          </button>
          <button
            onClick={() => setActiveTab('dates')}
            className={`px-4 py-2 rounded-lg transition-colors ${
              activeTab === 'dates'
                ? 'bg-indigo-500 text-white'
                : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
            }`}
          >
            Concert Dates
          </button>
        </div>

        <div className="bg-gray-800/50 rounded-lg p-6">
          {activeTab === 'info' && (
            <div className="space-y-4">
              <h2 className="text-2xl font-semibold mb-4">About</h2>
              <p className="text-gray-300">
                {artist.name} was formed in {artist.creationDate}. The band released their first album
                "{artist.firstAlbum}" and has since become a prominent name in the music industry.
              </p>
            </div>
          )}

          {activeTab === 'locations' && (
            <div>
              <h2 className="text-2xl font-semibold mb-4">Tour Locations</h2>
              <div className="grid md:grid-cols-2 gap-4">
                {artistLocations.map((location, index) => (
                  <div
                    key={index}
                    className="flex items-center space-x-2 bg-gray-700/30 p-3 rounded-lg"
                  >
                    <MapPin className="text-indigo-500" />
                    <span>{location}</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {activeTab === 'dates' && (
            <div>
              <h2 className="text-2xl font-semibold mb-4">Upcoming Concerts</h2>
              <div className="grid md:grid-cols-2 gap-4">
                {artistDates.map((date, index) => (
                  <div
                    key={index}
                    className="flex items-center space-x-2 bg-gray-700/30 p-3 rounded-lg"
                  >
                    <Calendar className="text-indigo-500" />
                    <span>{date}</span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ArtistDetails;