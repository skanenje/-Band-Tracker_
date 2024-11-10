export interface Artist {
  id: number;
  image: string;
  name: string;
  members: string[];
  creationDate: number;
  firstAlbum: string;
  locations: string;
  concertDates: string;
  relations: string;
}

export interface Location {
  id: number;
  name: string;
  locations: string[];
  dates: string;
}

export interface Date {
  id: number;
  name: string;
  dates: string[];
}

export interface Relation {
  id: number;
  name: string;
  datesLocations: Record<string, string[]>;
}