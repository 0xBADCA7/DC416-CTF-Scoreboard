module Mode.Scoreboard exposing (Scoreboard(..), view)

import Html exposing (..)
import Html.Attributes exposing (id)


type Scoreboard
    = Scoreboard (List Team)


type alias Team =
    { rank : Int
    , name : String
    , score : Int
    , lastSubmission : String
    }


view : Scoreboard -> Html msg
view (Scoreboard teams) =
    table [ id "scoreboard" ] <|
        List.append
            [ td [ id "scoreboardHeader" ]
                [ th [] [ text "Rank" ]
                , th [] [ text "Name" ]
                , th [] [ text "Score" ]
                , th [] [ text "Last Submission" ]
                ]
            ]
            (List.map viewTeam teams)


viewTeam : Team -> Html msg
viewTeam team =
    tr []
        [ td [] [ text (toString team.rank) ]
        , td [] [ text team.name ]
        , td [] [ text (toString team.score) ]
        , td [] [ text team.lastSubmission ]
        ]
