module Util exposing (prettyDate)

-- Core packages

import Date
import Time


-- HELPERS


prettyDate : Int -> String
prettyDate seconds =
    let
        asDate =
            seconds
                |> toFloat
                |> (\x -> x * Time.second)
                |> Date.fromTime

        dateStr =
            (toString <| Date.month asDate) ++ " " ++ (toString <| Date.day asDate)

        hourStr =
            asDate
                |> Date.hour
                |> toString
                |> String.padLeft 2 '0'

        minuteStr =
            asDate
                |> Date.minute
                |> toString
                |> String.padLeft 2 '0'
    in
        dateStr ++ " at " ++ hourStr ++ ":" ++ minuteStr
